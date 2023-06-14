import { BackendError, BackendNetworkError } from '../Services/BackendService';
// import { UserServiceErrors } from '../Services/UserService';


export const ErrorMap: Record<string, string> = {
//    ...UserServiceErrors
};


export class ErrorTranslator
{
    public static translate ( err: Error | BackendError | BackendNetworkError | null ): string | null
    {
        if ( ! err )
        {
            return null;
        }

        if ( BackendNetworkError.isBackendNetworkError(err) )
        {
            return 'Es konnte keine Verbindung zum Server hergestellt werden.\n' +
                '\n' +
                'Mögliche Ursachen sind:\n' +
                '- Dein PC / Tablet / Smartphone hat keine Internetverbindung\n' +
                '- Der Server ist aus technischen Gründen nicht erreichbar';
        }
        
        if ( BackendError.isBackendError(err) && ErrorMap[(err as BackendError).code] )
        {
            return ErrorMap[(err as BackendError).code] ;
        }

        return `Es ist ein Fehler aufgetreten: ${err.message}`;
    }
}
