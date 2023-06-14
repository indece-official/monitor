import { BackendService } from './BackendService';


export enum ConnectorV1Status
{
    Unregistered    = 'UNREGISTERED',
    Ready           = 'READY',
    Error           = 'ERROR'
}


export enum ConnectorV1CapabilitiesConfigValueType
{
    Text        = 'TEXT',
    Password    = 'PASSWORD',
    Number      = 'NUMBER',
    Select      = 'SELECT'
}


export interface ConnectorV1CapabilitiesConfigValue
{
    name:       string;
    label:      string;
    required:   boolean;
    type:       ConnectorV1CapabilitiesConfigValueType;
    options:    Array<string> | null;
}


export interface ConnectorV1Capabilities
{
    config_values:  Array<ConnectorV1CapabilitiesConfigValue>;
}


export interface ConnectorV1
{
    uid:            string;
    type:           string | null;
    host_uid:       string;
    version:        string | null;
    status:         ConnectorV1Status;
    connected:      boolean;
    error:          string | null;
    last_ping:      string | null;
}


export interface AddConnectorV1Request
{
    host_uid:   string;
}


export interface AddConnectorV1Result
{
    connector_uid:  string;
    config_file:    string;
}


export class ConnectorService
{
    private static _instance:           ConnectorService;
    private readonly _backendService:   BackendService;
    
    
    public static getInstance ( ): ConnectorService
    {
        if ( ! this._instance )
        {
            this._instance = new ConnectorService();
        }
        
        return this._instance;
    }


    constructor ( )
    {
        this._backendService = BackendService.getInstance();
    }


    public async getConnectors ( ): Promise<Array<ConnectorV1>>
    {
        const resp = await this._backendService.fetchJson(
            `/api/v1/connector`,
            {
                method: 'GET',
                headers:    {
                    'Accept':   'application/json'
                }
            }
        );

        return resp.connectors;
    }
    
    
    public async getConnector ( connectorUID: string ): Promise<ConnectorV1>
    {
        const resp = await this._backendService.fetchJson(
            `/api/v1/connector/${encodeURIComponent(connectorUID)}`,
            {
                method: 'GET',
                headers:    {
                    'Accept':   'application/json'
                }
            }
        );

        return resp.connector;
    }

    
    public async deleteConnector ( connectorUID: string ): Promise<void>
    {
        await this._backendService.fetchJson(
            `/api/v1/connector/${encodeURIComponent(connectorUID)}`,
            {
                method:     'DELETE',
                headers:    {
                    'Accept':   'application/json'
                }
            }
        );
    }


    public async addConnector ( params: AddConnectorV1Request ): Promise<AddConnectorV1Result>
    {
        const resp = await this._backendService.fetchJson(
            `/api/v1/connector`,
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
