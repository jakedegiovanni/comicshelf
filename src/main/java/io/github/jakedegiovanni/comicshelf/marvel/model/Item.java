package io.github.jakedegiovanni.comicshelf.marvel.model;

import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@NoArgsConstructor
public class Item {
    private String name;
    private String resourceURI;
}