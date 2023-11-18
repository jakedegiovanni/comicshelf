@import io.github.jakedegiovanni.comicshelf.views.Page
@import io.github.jakedegiovanni.comicshelf.marvel.model.DataWrapper
@import io.github.jakedegiovanni.comicshelf.marvel.model.Comic

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

                        <input type="hidden" name="name" value="${comic.getSeries().getName()}" />
                        <button type="submit" name="key" value="${resourceUri}" hx-post="/api/track">
                            <!-- todo following checks -->
                            <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg" style="width: 32px;height: 32px;"><g id="SVGRepo_bgCarrier" stroke-width="0"></g><g id="SVGRepo_tracerCarrier" stroke-linecap="round" stroke-linejoin="round"></g><g id="SVGRepo_iconCarrier"> <path fill-rule="evenodd" clip-rule="evenodd" d="M4 4C4 2.34315 5.34315 1 7 1H17C18.6569 1 20 2.34315 20 4V20.9425C20 22.6114 18.0766 23.5462 16.7644 22.5152L12 18.7717L7.23564 22.5152C5.92338 23.5462 4 22.6114 4 20.9425V4ZM7 3C6.44772 3 6 3.44772 6 4V20.9425L12 16.2283L18 20.9425V4C18 3.44772 17.5523 3 17 3H7Z" fill="#0F0F0F"></path> </g></svg>
                        </button>
                        <div class="pusher"></div>
                        <div class="marvel-links"><a href="series?series=${resourceUri}">${comic.getSeries().getName()}</a></div>
                        @for(var d : comic.getDates())
                            @if(d.getType().equals("onsaleDate"))
                                <div class="marvel-links">${d.getDate()}</div>
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
