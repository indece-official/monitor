import React from 'react';
import { ErrorBox } from '../../Components/ErrorBox/ErrorBox';
import { Spinner } from '../../Components/Spinner/Spinner';
import { ConfigService, ConfigPropertyV1 } from '../../Services/ConfigService';
import { List } from '../../Components/List/List';
import { ListEmpty } from '../../Components/List/ListEmpty';
import { ListItem } from '../../Components/List/ListItem';
import { ListItemHeader } from '../../Components/List/ListItemHeader';
import { ListItemHeaderField } from '../../Components/List/ListItemHeaderField';
import { PageContent } from '../../Components/PageContent/PageContent';
import { ListItemBody } from '../../Components/List/ListItemBody';


export interface ConfigPageProps
{
}


interface ConfigPageState
{
    config:     Array<ConfigPropertyV1>;
    loading:    boolean;
    error:      Error | null;
}


export class ConfigPage extends React.Component<ConfigPageProps, ConfigPageState>
{
    private readonly _configService: ConfigService;


    constructor ( props: ConfigPageProps )
    {
        super(props);

        this.state = {
            config:     [],
            loading:        false,
            error:          null
        };

        this._configService  = ConfigService.getInstance();
    }


    private async _load ( ): Promise<void>
    {
        try
        {
            this.setState({
                loading:    true,
                error:      null
            });

            const config = await this._configService.getConfig();

            this.setState({
                loading:    false,
                config:     Object.values(config)
            });
        }
        catch ( err )
        {
            console.error(`Error loading config: ${(err as Error).message}`, err);

            this.setState({
                loading:    false,
                error:      err as Error
            });
        }
    }


    public async componentDidMount ( ): Promise<void>
    {
        await this._load();
    }


    public render ( )
    {
        return (
            <PageContent>
                <h1>Configuration</h1>

                <ErrorBox error={this.state.error} />

                <List>
                    {this.state.config.length === 0 && !this.state.loading && !this.state.error ?
                        <ListEmpty>
                            No config properties found
                        </ListEmpty>
                    : null}

                    {this.state.config.map( ( configProperty ) => (
                        <ListItem key={configProperty.key}>
                            <ListItemHeader>
                                <ListItemHeaderField
                                    grow={true}
                                    text={configProperty.key}
                                />

                                {configProperty.value.length <= 100 ?
                                    <ListItemHeaderField
                                        text={configProperty.value}
                                    />
                                : null}
                            </ListItemHeader>

                            {configProperty.value.length > 100 ?
                                <ListItemBody>
                                    {configProperty.value}
                                </ListItemBody>
                            : null}
                        </ListItem>
                    ))}
                </List>

                <Spinner active={this.state.loading} />
            </PageContent>
        );
    }
}
