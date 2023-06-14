import { BackendService } from './BackendService';


export enum ConfigPropertyV1Key
{
    ConnectorHost   = 'CONNECTOR_HOST',
    ConnectorPort   = 'CONNECTOR_PORT',
    TlsCaCrt        = 'TLS_CA_CRT',
    TlsServerCrt    = 'TLS_SERVER_CRT'
}


export enum ConfigPropertyV1Bool
{
    True    = 'true',
    False   = 'false'
}


export interface ConfigPropertyV1
{
    key:        ConfigPropertyV1Key;
    value:      string;
    editable:   boolean;
}


export interface SetConfigPropertyV1Request
{
    key:    ConfigPropertyV1Key;
    value:  string;
}


export class ConfigService
{
    private static _instance:           ConfigService;
    private readonly _backendService:   BackendService;
    
    
    public static getInstance ( ): ConfigService
    {
        if ( ! this._instance )
        {
            this._instance = new ConfigService();
        }
        
        return this._instance;
    }


    constructor ( )
    {
        this._backendService = BackendService.getInstance();
    }


    public async getConfig ( ): Promise<Partial<Record<ConfigPropertyV1Key, ConfigPropertyV1>>>
    {
        const resp = await this._backendService.fetchJson(
            `/api/v1/config`,
            {
                method: 'GET',
                headers:    {
                    'Accept':       'application/json'
                }
            }
        );

        const config: Partial<Record<ConfigPropertyV1Key, ConfigPropertyV1>> = {};
        for ( const property of resp.properties )
        {
            config[property.key as ConfigPropertyV1Key] = property;
        }

        return config;
    }

    
    public async finishSetup ( ): Promise<void>
    {
        await this._backendService.fetchJson(
            `/api/v1/setup/finish`,
            {
                method:     'POST',
                headers:    {
                    'Accept':       'application/json',
                    'Content-Type': 'application/json'
                },
                body:   JSON.stringify({})
            }
        );
    }
   
   
    public async setConfigProperty ( key: ConfigPropertyV1Key, value: string ): Promise<void>
    {
        await this._backendService.fetchJson(
            `/api/v1/config/${encodeURIComponent(key)}`,
            {
                method:     'PUT',
                headers:    {
                    'Accept':       'application/json',
                    'Content-Type': 'application/json'
                },
                body:   JSON.stringify({
                    value
                })
            }
        );
    }
}
