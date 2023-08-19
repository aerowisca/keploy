import yaml
import os

def yaml_as_python(val):
    """Convert YAML to dict"""
    try:
        return yaml.safe_load_all(val)
    except yaml.YAMLError as exc:
        return exc

for file in os.listdir("../../../../keploy/mocks"):
    oldMocks = None
    newMocks = None
    with open('../../../../keploy/mocks/'+file,'r') as input_file:
        print("-------------------------------new file started")
        oldMocks = list(yaml_as_python(input_file))
    with open('../../../../keployTest990/mocks/'+file,'r') as input_file:
        newMocks = list(yaml_as_python(input_file))
        # print("the oldMocks are ",oldMocks)
    for valueOld,valueNew in zip(oldMocks,newMocks):
         # print("printing the grpc timeouts value: \n",valueOld["spec"]["grpcReq"]["headers"]["ordinary_headers"]["grpc-timeout"],valueNew["spec"]["grpcReq"]["headers"]["ordinary_headers"]["grpc-timeout"])
        if valueOld["spec"]["grpcReq"]["body"]==valueNew["spec"]["grpcReq"]["body"]:
            print("values matched")
        else:
            print("mistach found in file :%s",file)
            print("Want :%s\nGot :%s",valueOld["spec"]["grpcReq"]["body"],valueNew["spec"]["grpcReq"]["body"])
            break
        if valueOld["spec"]["grpcResp"]["body"]==valueNew["spec"]["grpcResp"]["body"]:
            print("values matched")
        else:
            print("mismatch found in file :%s",file)
            print("value mismatched \nWant :%s\nGot :%s",valueOld["spec"]["grpcResp"]["body"],valueNew["spec"]["grpcResp"]["body"])
            break
