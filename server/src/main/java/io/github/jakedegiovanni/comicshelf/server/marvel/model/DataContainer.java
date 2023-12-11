package io.github.jakedegiovanni.comicshelf.server.marvel.model;

import lombok.Data;
import lombok.NoArgsConstructor;

import java.util.List;

@Data @NoArgsConstructor
public class DataContainer<T> {
    private int offset;
    private int limit;
    private int total;
    private int count;
    private List<T> results;
}
