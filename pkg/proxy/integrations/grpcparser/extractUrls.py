import yaml
import os

urlfile = open('../pkg/proxy/integrations/grpcparser/apiUrl.txt', 'a')
for file in os.listdir("tests"):
    with open(os.path.join(os.getcwd(), "tests/" + file), 'r') as f:
        valuesYaml = yaml.load(f, Loader=yaml.FullLoader)
    url = valuesYaml['spec']['req']['url']
    urlfile.write(url + "\n")
urlfile.close()
