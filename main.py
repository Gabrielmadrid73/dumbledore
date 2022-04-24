from datetime import date, datetime
from kubernetes import client, config
from kubernetes.client.rest import ApiException
from base64 import b64encode
from fastapi import FastAPI
from pydantic import BaseModel
import json, boto3, botocore.exceptions, os

# KUBERNETES LOCAL CONFIG
# config.load_kube_config()

# KUBERNETES CLUSTER CONFIG
config.load_incluster_config()

# ENVIRONMENT VARIABLES
SECRET_API_VERSION = os.getenv('SECRET_API_VERSION', 'v1') 
AWS_REGION         = os.getenv('AWS_REGION', '')

# GLOBAL VARIABLES
APP         = FastAPI() # FASTAPI
API         = client.CoreV1Api() # KUBERNETES API
AWS_SSM     = boto3.client('ssm', region_name=AWS_REGION) # AWS SSM API


# CHECK IF SECRET EXIST
def check_secret(namespace="", secrets=[]):
    checked_secrets  = []
    secrets_notfound = []

    for secret in secrets:
        try:
            API.read_namespaced_secret(name=secret, namespace=namespace)
            checked_secrets.append(secret)
        except ApiException:
            secrets_notfound.append(secret)
            print("[WARNING] - Secret: " + secret + " NOT FOUND!!!")
    
    output = {
        "checked_secrets": checked_secrets,
        "secrets_notfound": secrets_notfound
    }

    check_return = json.dumps(output)

    return check_return

# FLOW TO UPDATE SECRET VALUE
def update_secret(namespace="", secrets=[]):

    secrets_updated                  = []
    secrets_error                    = []
    secrets_parameter_store_notfound = []
    
    # INITIATE THE COLLECTING DATA FROM CHECKED SECRETS
    for secret in secrets:
        # COLLECTING INFORMATIONS ABOUT EXISTING SECRET
        req_kubernetes             = API.read_namespaced_secret(name=secret, namespace=namespace).metadata.annotations
        req_kubernetes_json        = json.dumps(req_kubernetes, indent=4)
        req_kubernetes_json_values = json.loads(req_kubernetes_json)
        
        # DATA OF ANNOTATIONS
        # SETTING NAME OF PARAMETER STORE
        ssm_name = req_kubernetes_json_values["aws-ssm/aws-param-name"]
        # SETTING TYPE OF PARAMETER STORE
        ssm_type = req_kubernetes_json_values["aws-ssm/aws-param-type"]
        # SETTING KEY OF PARAMETER STORE
        ssm_key  = req_kubernetes_json_values["aws-ssm/aws-param-key"]

        # USED TO ALLOW THE OUTPUT GENERATED BY AWS
        def json_datetime_serializer(obj):
            if isinstance(obj, (datetime, date)):
                return obj.isoformat()
            raise TypeError("Type %s not serializable" % type(obj))

        ## COLLECTING VALUE OF PARAMETER STORE FROM AWS
        try:
            aws_return            = AWS_SSM.get_parameter(Name=ssm_name, WithDecryption=True)
            # CONVERTING RESPONSE TO JSON TO GET VALUE
            aws_return_parameters = json.dumps(aws_return['Parameter'], indent=4, default=json_datetime_serializer)
            aws_return_values     = json.loads(aws_return_parameters)   
            ssm_value             = aws_return_values["Value"]

            try:
                # ENCODING PARAMETER STORE VALUE. IT'S REQUIRED BY KUBERNETES API SERVER
                ssm_value_encoded = b64encode(ssm_value.encode("utf-8")).decode()
                
                # MANIFEST TO BE SEND TO KUBERNETES API SERVER
                body = {
                    "apiVersion": SECRET_API_VERSION,
                    "kind": "Secret",
                    "data": { ssm_type : ssm_value_encoded },
                    "metadata": {"annotations": {"aws-ssm/aws-param-key": ssm_key, "aws-ssm/aws-param-name": ssm_name, "aws-ssm/aws-param-type": ssm_type }, "name": secret, "namespace": namespace},
                    "type": "Opaque"
                }
                
                # REQUEST TO UPDATE SECRET
                API.patch_namespaced_secret(secret, namespace, body=body)
                
                # ADD SECRET TO UPDATED LIST
                secrets_updated.append(secret)
                
                # STDOUT OF UPDATES
                print("[INFO] Secret: " + secret + " - updated with value from Parameter Store: " + ssm_name)
                
            # EXCEPTION TO A ERROR FROM REQUEST TO AWS
            except ApiException as e:
                # ADD SECRET TO A LIST OF SECRETS WITH ERROR
                secrets_error.append(secret)
                print("Exception when calling CoreV1Api->patch_namespaced_secret: %s\n" % e)
                
        # EXCEPTION TO A INEXISTING PARAMETER STORE
        except botocore.exceptions.ClientError as error:
            if error.response['Error']['Code'] == 'ParameterNotFound':                
                # ADD SECRET TO A LIST OF PARAMETER STORE NOT FOUND
                secrets_parameter_store_notfound.append(secret)
                print("[WARNING] - Parameter store: " + ssm_name + " NOT FOUND!! Referred of secret: " + secret)
            else:
                raise error
        
    output = {
        "secrets_updated": secrets_updated,
        "secrets_error": secrets_error,
        "secrets_parameter_store_notfound": secrets_parameter_store_notfound
    }

    update_secret_return = json.dumps(output)

    return update_secret_return

def backend(namespace="", secrets=[]):
    check           = check_secret(namespace, secrets)
    checked_return  = json.loads(check)
    checked_secrets = checked_return["checked_secrets"]

    sync_secrets        = update_secret(namespace, checked_secrets)
    sync_secrets_return = json.loads(sync_secrets)

    secrets_updated                  = sync_secrets_return["secrets_updated"]
    secrets_error                    = sync_secrets_return["secrets_error"]
    secrets_parameter_store_notfound = sync_secrets_return["secrets_parameter_store_notfound"]
    secrets_notfound                 = checked_return["secrets_notfound"]
    
    api_response = {
        "secrets_updated": secrets_updated,
        "secrets_notfound": secrets_notfound,
        "secrets_error": secrets_error,
        "secrets_parameter_store_notfound": secrets_parameter_store_notfound
    }
    
    return api_response

# VALIDATE SCHEMA OF PAYLOAD
class schema(BaseModel):
    namespace: str
    secrets: list

# FAST API
@APP.post("/api/v1/secrets/sync")
async def request(payload: schema):
    request      = payload.dict()
    namespace    = request["namespace"]
    secrets      = request["secrets"]
    api_response = backend(namespace, secrets)

    return api_response

# HEALTH CHECK
@APP.get("/health-check")
async def api_request():
    return "UP"