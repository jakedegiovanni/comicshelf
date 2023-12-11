package io.github.jakedegiovanni.comicshelf.sdk.marvel.model;

import com.fasterxml.jackson.annotation.JsonIgnore;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.util.List;

@Data @NoArgsConstructor
public class Comic {
    private int id;
    private String title;
    private String resourceURI;
    private List<Url> urls;
    private String modified;
    private Thumbnail thumbnail;
    private String format;
    private int issuerNumber;
    private Item series;
    private List<Date> dates;
    @JsonIgnore private boolean following;
}
