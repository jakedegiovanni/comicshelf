@import io.github.jakedegiovanni.comicshelf.server.views.Page
@import io.github.jakedegiovanni.comicshelf.server.marvel.model.DataWrapper
@import io.github.jakedegiovanni.comicshelf.server.marvel.model.Comic
@import java.time.format.DateTimeFormatter

@param Page<DataWrapper<Comic>> page

@template.index(
    page = page,
    content = @`
        <base href="/marvel-unlimited/" />

        @for(var comic : page.getContents().getData().getResults())
            @if(!comic.getTitle().toLowerCase().contains("variant"))
                !{var resourceUri = comic.getSeries().getResourceURI();}
                <div id="${comic.getId()}" class="container">
                    <img src="${comic.getThumbnail().getPath()}/portrait_uncanny.${comic.getThumbnail().getExtension()}" alt="img"/>
                    <form>
                        <div class="title"><h3>${comic.getTitle()}</h3></div>
                        <input type="hidden" name="internalId" value="${comic.getId()}" />
                        <input type="hidden" name="title" value="${comic.getTitle()}" />
                        <input type="hidden" name="series" value="${comic.getSeries().getName()}" />
                        @if(comic.isFollowing())
                            @template.unfollow()
                        @else
                            @template.follow()
                        @endif
                        <div class="pusher"></div>
                        <div class="marvel-links"><a href="series?series=${resourceUri}">${comic.getSeries().getName()}</a></div>
                        @for(var d : comic.getDates())
                            @if(d.getType().equals("onsaleDate"))
                                <div class="marvel-links">${d.getDate().format(DateTimeFormatter.ISO_LOCAL_DATE)}</div>
                            @endif
                        @endfor
                        @for(var u : comic.getUrls())
                            @if(u.getType().equals("detail"))
                                <div class="marvel-links"><a href="${u.getUrl()}" target="_blank">Details</a></div>
                            @endif
                        @endfor
                        @for(var u : comic.getUrls())
                            @if(u.getType().equals("reader"))
                                <div class="marvel-links"><a href="${u.getUrl()}" target="_blank">Reader</a></div>
                            @endif
                        @endfor
                    </form>
                </div>
            @endif
        @endfor
    `,
    footer = @`<a href="https://marvel.com">${page.getContents().getAttributionText()}</a>`
)
