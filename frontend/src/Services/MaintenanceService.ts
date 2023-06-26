import { BackendService } from './BackendService';


export interface MaintenanceV1Affected
{
    host_uids:  Array<string>;
    check_uids: Array<string>;
    tag_uids:   Array<string>;
}


export interface MaintenanceV1
{
    uid:                string;
    title:              string;
    message:            string;
    affected:           MaintenanceV1Affected;
    datetime_created:   string;
    datetime_updated:   string;
    datetime_start:     string;
    datetime_finish:    string | null;
}


export interface AddMaintenanceV1Request
{
    title:              string;
    message:            string;
    affected:           MaintenanceV1Affected;
    datetime_start:     string;
    datetime_finish:    string | null;
}


export interface UpdateMaintenanceV1Request
{
    title:              string;
    message:            string;
    affected:           MaintenanceV1Affected;
    datetime_start:     string;
    datetime_finish:    string | null;
}


export class MaintenanceService
{
    private static _instance:           MaintenanceService;
    private readonly _backendService:   BackendService;
    
    
    public static getInstance ( ): MaintenanceService
    {
        if ( ! this._instance )
        {
            this._instance = new MaintenanceService();
        }
        
        return this._instance;
    }


    constructor ( )
    {
        this._backendService = BackendService.getInstance();
    }


    public async getMaintenances ( active: boolean | null, from: number | null, size: number | null ): Promise<Array<MaintenanceV1>>
    {
        const resp = await this._backendService.fetchJson(
            `/api/v1/maintenance?active=${encodeURIComponent(active ?? '')}&from=${encodeURIComponent(from ?? '')}&size=${encodeURIComponent(size ?? '')}`,
            {
                method: 'GET',
                headers:    {
                    'Accept':   'application/json'
                }
            }
        );

        return resp.maintenances;
    }
    
    
    public async getMaintenance ( maintenanceUID: string ): Promise<MaintenanceV1>
    {
        const resp = await this._backendService.fetchJson(
            `/api/v1/maintenance/${encodeURIComponent(maintenanceUID)}`,
            {
                method: 'GET',
                headers:    {
                    'Accept':   'application/json'
                }
            }
        );

        return resp.maintenance;
    }

    
    public async deleteMaintenance ( maintenanceUID: string ): Promise<void>
    {
        await this._backendService.fetchJson(
            `/api/v1/maintenance/${encodeURIComponent(maintenanceUID)}`,
            {
                method:     'DELETE',
                headers:    {
                    'Accept':   'application/json'
                }
            }
        );
    }


    public async addMaintenance ( params: AddMaintenanceV1Request ): Promise<string>
    {
        const resp = await this._backendService.fetchJson(
            `/api/v1/maintenance`,
            {
                method: 'POST',
                headers:    {
                    'Accept':       'application/json',
                    'Content-Type': 'application/json'
                },
                body:   JSON.stringify(params)
            }
        );

        return resp.maintenance_uid;
    }


    public async updateMaintenance ( maintenanceUID: string, params: UpdateMaintenanceV1Request ): Promise<void>
    {
        await this._backendService.fetchJson(
            `/api/v1/maintenance/${encodeURIComponent(maintenanceUID)}`,
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
}
