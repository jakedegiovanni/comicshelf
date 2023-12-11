package io.github.jakedegiovanni.comicshelf.server.views;

import lombok.AllArgsConstructor;
import lombok.Data;

import java.time.LocalDate;
import java.time.format.DateTimeFormatter;

@Data @AllArgsConstructor
public class Index {

    private final LocalDate date;

    public String getDate() {
        return date.format(DateTimeFormatter.ISO_LOCAL_DATE);
    }
}
