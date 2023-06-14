// import { StateSubject } from 'ts-subject';


export class ModalService
{
    private static _instance:               ModalService;
    // private readonly _subjectDeleteDataset: StateSubject<DeleteDatasetModalParams | null>;
    

    public static getInstance ( ): ModalService
    {
        if ( ! this._instance )
        {
            this._instance = new ModalService();
        }

        return this._instance;
    }


    /* constructor ( )
    {
        this._subjectDeleteDataset = new StateSubject<DeleteDatasetModalParams | null>(null);
    }
   
   
    public showDeleteDataset ( params: DeleteDatasetModalParams ): void
    {
        this._subjectDeleteDataset.next(params);
    }
    
    
    public hideDeleteDataset ( ): void
    {
        this._subjectDeleteDataset.next(null);
    }


    public getDeleteDataset ( ): StateSubject<DeleteDatasetModalParams | null>
    {
        return this._subjectDeleteDataset;
    }*/
}
