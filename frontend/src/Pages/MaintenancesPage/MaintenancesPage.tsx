import React from 'react';
import DayJS from 'dayjs';
import { ErrorBox } from '../../Components/ErrorBox/ErrorBox';
import { Spinner } from '../../Components/Spinner/Spinner';
import { MaintenanceService, MaintenanceV1 } from '../../Services/MaintenanceService';
import { Button } from '../../Components/Button/Button';
import { LinkUtils } from '../../utils/LinkUtils';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faPlus, faTrash } from '@fortawesome/free-solid-svg-icons';
import { List } from '../../Components/List/List';
import { ListEmpty } from '../../Components/List/ListEmpty';
import { ListItem } from '../../Components/List/ListItem';
import { ListItemHeaderField } from '../../Components/List/ListItemHeaderField';
import { ListItemHeader } from '../../Components/List/ListItemHeader';
import { ListItemHeaderAction } from '../../Components/List/ListItemHeaderAction';
import { TagService, TagV1 } from '../../Services/TagService';
import { HostService, HostV1 } from '../../Services/HostService';
import { Tag } from '../../Components/Tag/Tag';


export interface MaintenancesPageProps
{
}


interface MaintenancesPageState
{
    maintenances:   Array<MaintenanceV1>;
    hostsMap:       Record<string, HostV1>;
    tagsMap:        Record<string, TagV1>;
    loading:        boolean;
    error:          Error | null;
}


export class MaintenancesPage extends React.Component<MaintenancesPageProps, MaintenancesPageState>
{
    private readonly _maintenanceService: MaintenanceService;
    private readonly _hostService: HostService;
    private readonly _tagService: TagService;


    constructor ( props: MaintenancesPageProps )
    {
        super(props);

        this.state = {
            maintenances:   [],
            hostsMap:       {},
            tagsMap:        {},
            loading:        false,
            error:          null
        };

        this._maintenanceService  = MaintenanceService.getInstance();
        this._hostService  = HostService.getInstance();
        this._tagService  = TagService.getInstance();
    }


    private async _load ( ): Promise<void>
    {
        try
        {
            this.setState({
                loading:    true,
                error:      null
            });

            const maintenances = await this._maintenanceService.getMaintenances(false, 0, 50);  // TODO
            const hosts = await this._hostService.getHosts();
            const hostsMap: Record<string, HostV1> = {};
            for ( const host of hosts )
            {
                hostsMap[host.uid] = host;
            }

            const tags = await this._tagService.getTags();
            const tagsMap: Record<string, TagV1> = {};
            for ( const tag of tags )
            {
                tagsMap[tag.uid] = tag;
            }

            this.setState({
                loading:    false,
                maintenances,
                hostsMap,
                tagsMap
            });
        }
        catch ( err )
        {
            console.error(`Error loading maintenances: ${(err as Error).message}`, err);

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
            <div className='MaintenancesPage'>
                <h1>Maintenances</h1>

                <ErrorBox error={this.state.error} />

                <Button to={LinkUtils.make('maintenance', 'add')}>
                    <FontAwesomeIcon icon={faPlus} /> Add a maintenance
                </Button>

                <List>
                    {this.state.maintenances.length === 0 && !this.state.loading && !this.state.error ?
                        <ListEmpty>
                            No maintenances found
                        </ListEmpty>
                    : null}

                    {this.state.maintenances.map( ( maintenance ) => (
                        <ListItem key={maintenance.uid}>
                            <ListItemHeader>
                                <ListItemHeaderField
                                    to={LinkUtils.make('maintenance', maintenance.uid)}
                                    grow={true}
                                    text={maintenance.title}
                                    subtext={`${DayJS(maintenance.datetime_start).format()}${maintenance.datetime_finish ? ' - ' + DayJS(maintenance.datetime_finish).format() : ''}`}
                                />

                                <ListItemHeaderField
                                    to={LinkUtils.make('maintenance', maintenance.uid)}>
                                    {maintenance.affected.host_uids.map( ( hostUID ) => (
                                        <div key={hostUID}>
                                            {this.state.hostsMap[hostUID]?.name || '?'}
                                        </div>
                                    ))}

                                    {maintenance.affected.tag_uids.map( ( tagUID ) => (
                                        <Tag key={tagUID} color={this.state.tagsMap[tagUID]?.color || ''}>
                                            {this.state.tagsMap[tagUID]?.name || '?'}
                                        </Tag>
                                    ))}
                                </ListItemHeaderField>

                                <ListItemHeaderAction
                                    to={LinkUtils.make('maintenance', maintenance.uid, 'delete')}
                                    icon={faTrash}
                                />
                            </ListItemHeader>
                        </ListItem>
                    ))}
                </List>

                <Spinner active={this.state.loading} />
            </div>
        );
    }
}
