import { BackendService } from './BackendService';


export interface TagV1
{
    uid:    string;
    name:   string;
    color:  string;
}


export interface AddTagV1Request
{
    name:   string;
    color:  string;
}


export interface UpdateTagV1Request
{
    name:   string;
    color:  string;
}


export class TagService
{
    private static _instance:           TagService;
    private readonly _backendService:   BackendService;
    
    
    public static getInstance ( ): TagService
    {
        if ( ! this._instance )
        {
            this._instance = new TagService();
        }
        
        return this._instance;
    }


    constructor ( )
    {
        this._backendService = BackendService.getInstance();
    }


    public async getTags ( ): Promise<Array<TagV1>>
    {
        const resp = await this._backendService.fetchJson(
            `/api/v1/tag`,
            {
                method: 'GET',
                headers:    {
                    'Accept':       'application/json'
                }
            }
        );

        return resp.tags;
    }
    
    
    public async getTag ( tagUID: string ): Promise<TagV1>
    {
        const resp = await this._backendService.fetchJson(
            `/api/v1/tag/${encodeURIComponent(tagUID)}`,
            {
                method: 'GET',
                headers:    {
                    'Accept':       'application/json'
                }
            }
        );

        return resp.tag;
    }


    public async addTag ( params: AddTagV1Request ): Promise<string>
    {
        const resp = await this._backendService.fetchJson(
            `/api/v1/tag`,
            {
                method: 'POST',
                headers:    {
                    'Accept':       'application/json',
                    'Content-Type': 'application/json'
                },
                body:   JSON.stringify(params)
            }
        );

        return resp.tag_uid;
    }


    public async updateTag ( tagUID: string, params: UpdateTagV1Request ): Promise<void>
    {
        await this._backendService.fetchJson(
            `/api/v1/tag/${encodeURIComponent(tagUID)}`,
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
    

    public async deleteTag ( tagUID: string ): Promise<void>
    {
        await this._backendService.fetchJson(
            `/api/v1/tag/${encodeURIComponent(tagUID)}`,
            {
                method: 'DELETE',
                headers:    {
                    'Accept':       'application/json'
                }
            }
        );
    }
}
