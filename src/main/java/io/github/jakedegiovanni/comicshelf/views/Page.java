package io.github.jakedegiovanni.comicshelf.views;

import lombok.EqualsAndHashCode;
import lombok.Getter;
import lombok.ToString;
import org.springframework.ui.Model;

import java.time.LocalDate;

@Getter
@ToString(callSuper = true)
@EqualsAndHashCode(callSuper = true)
public class Page<T> extends Index {

    private final T data;

    public Page(LocalDate date, T data) {
        super(date);
        this.data = data;
    }

    // todo - happy depending on spring here?
    public static <T> void setupModel(Model model, T data, LocalDate date) {
        model.addAttribute("page", new Page<T>(date, data));
    }
}
