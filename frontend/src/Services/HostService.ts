import { BackendService } from './BackendService';
import { TagV1 } from './TagService';


export interface HostV1Status
{
    count_critical: number;
    count_warning:  number;
    count_ok:       number;
    count_unknown:  number;
}


export interface HostV1
{
    uid:    string;
    name:   string;
    tags:   Array<TagV1>;
    status: HostV1Status;
}


export interface AddHostV1Request
{
    name:       string;
    tag_uids:   Array<string>;
}


export interface UpdateHostV1Request
{
    name:       string;
    tag_uids:   Array<string>;
}


export class HostService
{
    private static _instance:           HostService;
    private readonly _backendService:   BackendService;
    
    
    public static getInstance ( ): HostService
    {
        if ( ! this._instance )
        {
            this._instance = new HostService();
        }
        
        return this._instance;
    }


    constructor ( )
    {
        this._backendService = BackendService.getInstance();
    }


    public async getHosts ( ): Promise<Array<HostV1>>
    {
        const resp = await this._backendService.fetchJson(
            `/api/v1/host`,
            {
                method: 'GET',
                headers:    {
                    'Accept':       'application/json'
                }
            }
        );

        return resp.hosts;
    }
    
    
    public async getHost ( hostUID: string ): Promise<HostV1>
    {
        const resp = await this._backendService.fetchJson(
            `/api/v1/host/${encodeURIComponent(hostUID)}`,
            {
                method: 'GET',
                headers:    {
                    'Accept':       'application/json'
                }
            }
        );

        return resp.host;
    }


    public async addHost ( params: AddHostV1Request ): Promise<string>
    {
        const resp = await this._backendService.fetchJson(
            `/api/v1/host`,
            {
                method: 'POST',
                headers:    {
                    'Accept':       'application/json',
                    'Content-Type': 'application/json'
                },
                body:   JSON.stringify(params)
            }
        );

        return resp.host_uid;
    }
    
    
    public async updateHost ( hostUID: string, params: UpdateHostV1Request ): Promise<void>
    {
        await this._backendService.fetchJson(
            `/api/v1/host/${encodeURIComponent(hostUID)}`,
            {
                method: 'PUT',
                headers:    {
                    'Accept':       'application/json',
                    'Content-Type': 'application/json'
                },
                body:   JSON.stringify(params)
            }
        );
    }
    

    public async deleteHost ( hostUID: string ): Promise<void>
    {
        await this._backendService.fetchJson(
            `/api/v1/host/${encodeURIComponent(hostUID)}`,
            {
                method: 'DELETE',
                headers:    {
                    'Accept':       'application/json'
                }
            }
        );
    }
}
