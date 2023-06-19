import DayJS from 'dayjs';
import { CheckStatusV1Status } from '../Services/CheckService';
import { AgentV1, AgentV1Status } from '../Services/AgentService';
import { UserV1Role } from '../Services/UserService';


const AGENT_STATUS_MAP: Record<AgentV1Status, string> = {
    [AgentV1Status.Unregistered]:   'Configuring ...',
    [AgentV1Status.Ready]:          'Ready',
    [AgentV1Status.Error]:          'Error',
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


    public static agentStatus ( agent: AgentV1 ): string
    {
        return `${AGENT_STATUS_MAP[agent.status] || agent.status} (${agent.connected ? 'Online' : 'Offline'})`;
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