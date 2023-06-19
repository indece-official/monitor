import { BackendService } from './BackendService';


export enum CheckerV1ParamType
{
    Text        = 'TEXT',
    Password    = 'PASSWORD',
    Number      = 'NUMBER',
    Select      = 'SELECT',
    Duration    = 'DURATION',
    Boolean     = 'BOOLEAN'
}


export interface CheckerV1Param
{
    name:       string;
    label:      string;
    hint:       string | null;
    required:   boolean;
    type:       CheckerV1ParamType;
    options:    Array<string>;
}


export interface CheckerV1Capabilities
{
    params:             Array<CheckerV1Param>;
    default_schedule:   string | null;
}


export interface CheckerV1
{
    uid:            string;
    name:           string;
    type:           string;
    agent_type:     string;
    version:        string;
    custom_checks:  boolean;
    capabilities:   CheckerV1Capabilities;
}


export class CheckerService
{
    private static _instance:           CheckerService;
    private readonly _backendService:   BackendService;
    
    
    public static getInstance ( ): CheckerService
    {
        if ( ! this._instance )
        {
            this._instance = new CheckerService();
        }
        
        return this._instance;
    }


    constructor ( )
    {
        this._backendService = BackendService.getInstance();
    }


    public async getCheckers ( ): Promise<Array<CheckerV1>>
    {
        const resp = await this._backendService.fetchJson(
            `/api/v1/checker`,
            {
                method: 'GET',
                headers:    {
                    'Accept':       'application/json'
                }
            }
        );

        return resp.checkers;
    }


    public async getChecker ( checkerUID: string ): Promise<CheckerV1>
    {
        const resp = await this._backendService.fetchJson(
            `/api/v1/checker/${encodeURIComponent(checkerUID)}`,
            {
                method: 'GET',
                headers:    {
                    'Accept':       'application/json'
                }
            }
        );

        return resp.checker;
    }
}
