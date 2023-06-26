export class ListUtils
{
    public static intersects<T> ( listA: Array<T>, listB: Array<T> ): boolean
    {
        for ( const elemA of listA )
        {
            for ( const elemB of listB )
            {
                if ( elemA === elemB )
                {
                    return true;
                }
            }
        }

        return false;
    }
}
