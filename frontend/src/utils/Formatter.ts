import DayJS from 'dayjs';
import { CheckStatusV1Status } from '../Services/CheckService';
import { ConnectorV1, ConnectorV1Status } from '../Services/ConnectorService';
import { UserV1Role } from '../Services/UserService';


const CONNECTOR_STATUS_MAP: Record<ConnectorV1Status, string> = {
    [ConnectorV1Status.Unregistered]:   'Configuring ...',
    [ConnectorV1Status.Ready]:          'Ready',
    [ConnectorV1Status.Error]:          'Error',
};


const CHECKSTATUS_STATUS_MAP: Record<CheckStatusV1Status, string> = {
    [CheckStatusV1Status.Critical]: 'CRIT',
    [CheckStatusV1Status.Warning]:  'WARN',
    [CheckStatusV1Status.Ok]:       'OK',
    [CheckStatusV1Status.Unknown]:  'UNKN',
};


const USER_ROLE_MAP: Record<UserV1Role, string> = {
    [UserV1Role.Show]:  'Show',
    [UserV1Role.Admin]: 'Admin'
};


export class Formatter
{
    public static datetime ( str: string | DayJS.Dayjs | null ): string
    {
        if ( ! str )
        {
            return '-';
        }

        return DayJS(str).format('DD.MM.YYYY HH:mm') + ' Uhr';
    }


    public static connectorStatus ( connector: ConnectorV1 ): string
    {
        return `${CONNECTOR_STATUS_MAP[connector.status] || connector.status} (${connector.connected ? 'Online' : 'Offline'})`;
    }
    
    
    public static checkStatus ( checkStatus: CheckStatusV1Status ): string
    {
        return CHECKSTATUS_STATUS_MAP[checkStatus] || checkStatus;
    }


    public static userRole ( userRole: UserV1Role ): string
    {
        return USER_ROLE_MAP[userRole] || userRole;
    }
}