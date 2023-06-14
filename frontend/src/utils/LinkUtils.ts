export class LinkUtils
{
    public static make ( ...segments: Array<string | number | boolean> ): string
    {
        return '/' + segments
            .map(encodeURIComponent)
            .join('/');
    }
}
