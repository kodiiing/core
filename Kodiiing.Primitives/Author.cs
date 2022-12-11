namespace Kodiiing.Primitives
{
    public class Author
    {
        public readonly string Name;
        public readonly  Uri ProfileUrl;
        public readonly Uri? PictureUri;

        public Author(string name, Uri profileUrl, Uri? pictureUri = null)
        {
            Name = name ?? throw new ArgumentNullException(nameof(name));
            ProfileUrl = profileUrl ?? throw new ArgumentNullException(nameof(profileUrl));
            PictureUri = pictureUri;
        }
    };
}