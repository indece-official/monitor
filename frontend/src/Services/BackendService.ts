import Cookies from 'universal-cookie';
import { Environment } from '../utils/Environment';


export class BackendNetworkError extends Error
{
    public readonly url:        string;


    constructor ( url: string,
                  message: string )
    {
        super(`Network error sending request to server: '${url}' => ${message}`);

        this.url        = url;
    }


    public static isBackendNetworkError ( err: Error )
    {
        return err instanceof BackendNetworkError;
    }
}


export class BackendError extends Error
{
    public readonly url:        string;
    public readonly status:     number;
    public readonly statusText: string;
    public readonly code:       string;
    public readonly details:    string;


    constructor ( url: string,
                  status: number,
                  statusText: string,
                  code?: string,
                  details?: string )
    {
        super(`Error sending request to server: '${url}' => ${status} ${statusText} (${code || '-'}: ${details || '-'})`);

        this.url        = url;
        this.status     = status;
        this.statusText = statusText;
        this.code       = code || '';
        this.details    = details || '';
    }


    public static isBackendError ( err: Error )
    {
        return err instanceof BackendError;
    }
}


export class BackendService
{
    private static _instance: BackendService;
    
    
    public static getInstance ( ): BackendService
    {
        if ( ! this._instance )
        {
            this._instance = new BackendService();
        }
        
        return this._instance;
    }
    
    
    private readonly _cookies: Cookies;
    private _accessToken: string | null;


    constructor ( )
    {
        this._cookies = new Cookies();
        this._accessToken = null;
    }


    private async _fetchJson ( url: string, init?: RequestInit ): Promise<any>
    {
        let resp: Response;

        init = {
            ...init,
            headers: {
                ...(init?.headers || {})
            }
        };

        if ( Environment.setup.token )
        {
            (init.headers as Record<string, string>)['X-Setup-Token'] = (init.headers as Record<string, string>)['X-Setup-Token'] || Environment.setup.token;
        }
        
        if ( this._accessToken )
        {
            (init.headers as Record<string, string>).Authorization = (init.headers as Record<string, string>).Authorization || `Bearer ${this._accessToken}`;
        }

        try
        {
            resp = await fetch(url, {
                credentials: 'include',
                ...init
            });
        }
        catch ( err )
        {
            throw new BackendNetworkError(url, (err as Error).message);
        }

        if ( ! resp.ok )
        {
            if ( resp.status >= 400 )
            {
                let body: any = null;

                try
                {
                    body = await resp.json();
                }
                catch ( err ) { }

                if ( body )
                {
                    throw new BackendError(
                        url,
                        resp.status,
                        resp.statusText,
                        body.error.code,
                        body.error.details
                    );
                }
            }

            throw new BackendError(
                url,
                resp.status,
                resp.statusText
            );
        }

        try
        {
            return await resp.json();
        }
        catch ( err )
        {
            throw new BackendNetworkError(url, (err as Error).message);
        }
    }


    private async _fetchText ( url: string, init?: RequestInit ): Promise<string>
    {
        let resp: Response;

        if ( this._accessToken )
        {
            init = {
                ...init,
                headers: {
                    Authorization: `Bearer ${this._accessToken}`,
                    ...(init?.headers || {})
                }
            };
        }

        try
        {
            resp = await fetch(url, {
                credentials: 'include',
                ...init
            });
        }
        catch ( err )
        {
            throw new BackendNetworkError(url, (err as Error).message);
        }

        if ( ! resp.ok )
        {
            if ( resp.status >= 400 )
            {
                let body: any = null;

                try
                {
                    body = await resp.text();
                }
                catch ( err ) { }

                if ( body )
                {
                    throw new BackendError(
                        url,
                        resp.status,
                        resp.statusText,
                        body.error.code,
                        body.error.details
                    );
                }
            }

            throw new BackendError(
                url,
                resp.status,
                resp.statusText
            );
        }

        try
        {
            return await resp.text();
        }
        catch ( err )
        {
            throw new BackendNetworkError(url, (err as Error).message);
        }
    }

    
    private async _getCsrfToken ( ): Promise<string>
    {
        await this._fetchText(
            `/api/user/v1/csrf`,
            {
                method:     'HEAD',
                headers:    {
                    'Accept':       'application/json'
                }
            }
        );

        return this._cookies.get('csrf');
    }


    public async fetchJson ( url: string, init?: RequestInit ): Promise<any>
    {
        let csrfToken = '';

        if ( init && init.method && ['POST', 'PUT', 'DELETE', 'PATCH'].includes(init.method.trim().toUpperCase()) )
        {
            try
            {
                csrfToken = await this._getCsrfToken();
            }
            catch ( err )
            {
                console.error(`Couldn't load csrf token: ${(err as Error).message}`, err);
            }
        }

        init = init || {};
        init.headers = init.headers || {};

        if ( csrfToken )
        {
            (init.headers as any)['X-CSRF-Token'] = csrfToken;
        }

        return this._fetchJson(url, init);
    }


    public setAccessToken ( accessToken: string ): void
    {
        this._accessToken = accessToken;
    }


    public deleteAccessToken ( ): void
    {
        this._accessToken = null;
    }
}
