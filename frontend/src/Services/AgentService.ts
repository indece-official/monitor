import { BackendService } from './BackendService';


export enum AgentV1Status
{
    Unregistered    = 'UNREGISTERED',
    Ready           = 'READY',
    Error           = 'ERROR'
}


export enum AgentV1CapabilitiesConfigValueType
{
    Text        = 'TEXT',
    Password    = 'PASSWORD',
    Number      = 'NUMBER',
    Select      = 'SELECT'
}


export interface AgentV1CapabilitiesConfigValue
{
    name:       string;
    label:      string;
    required:   boolean;
    type:       AgentV1CapabilitiesConfigValueType;
    options:    Array<string> | null;
}


export interface AgentV1Capabilities
{
    config_values:  Array<AgentV1CapabilitiesConfigValue>;
}


export interface AgentV1
{
    uid:            string;
    type:           string | null;
    host_uid:       string;
    version:        string | null;
    status:         AgentV1Status;
    connected:      boolean;
    error:          string | null;
    last_ping:      string | null;
}


export interface AddAgentV1Request
{
    host_uid:   string;
}


export interface AddAgentV1Result
{
    agent_uid:  string;
    config_file:    string;
}


export class AgentService
{
    private static _instance:           AgentService;
    private readonly _backendService:   BackendService;
    
    
    public static getInstance ( ): AgentService
    {
        if ( ! this._instance )
        {
            this._instance = new AgentService();
        }
        
        return this._instance;
    }


    constructor ( )
    {
        this._backendService = BackendService.getInstance();
    }


    public async getAgents ( ): Promise<Array<AgentV1>>
    {
        const resp = await this._backendService.fetchJson(
            `/api/v1/agent`,
            {
                method: 'GET',
                headers:    {
                    'Accept':   'application/json'
                }
            }
        );

        return resp.agents;
    }
    
    
    public async getAgent ( agentUID: string ): Promise<AgentV1>
    {
        const resp = await this._backendService.fetchJson(
            `/api/v1/agent/${encodeURIComponent(agentUID)}`,
            {
                method: 'GET',
                headers:    {
                    'Accept':   'application/json'
                }
            }
        );

        return resp.agent;
    }

    
    public async deleteAgent ( agentUID: string ): Promise<void>
    {
        await this._backendService.fetchJson(
            `/api/v1/agent/${encodeURIComponent(agentUID)}`,
            {
                method:     'DELETE',
                headers:    {
                    'Accept':   'application/json'
                }
            }
        );
    }


    public async addAgent ( params: AddAgentV1Request ): Promise<AddAgentV1Result>
    {
        const resp = await this._backendService.fetchJson(
            `/api/v1/agent`,
            {
                method: 'POST',
                headers:    {
                    'Accept':       'application/json',
                    'Content-Type': 'application/json'
                },
                body:   JSON.stringify(params)
            }
        );

        return resp;
    }
}
