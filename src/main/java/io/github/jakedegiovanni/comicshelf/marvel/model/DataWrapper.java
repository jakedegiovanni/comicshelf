package io.github.jakedegiovanni.comicshelf.marvel.model;

import lombok.Data;
import lombok.NoArgsConstructor;

@Data @NoArgsConstructor
public class DataWrapper<T> {
    private int code;
    private String status;
    private String copyright;
    private String attributionText;
    private String attributionHTML;
    private String etag;
    private DataContainer<T> data;
}
