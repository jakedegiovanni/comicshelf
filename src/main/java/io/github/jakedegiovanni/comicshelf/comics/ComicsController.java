package io.github.jakedegiovanni.comicshelf.comics;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.Data;
import lombok.NoArgsConstructor;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Controller;
import org.springframework.ui.Model;
import org.springframework.web.bind.annotation.ModelAttribute;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;

@Controller
@RequestMapping("/api/comics")
@RequiredArgsConstructor
public class ComicsController {

    public record ComicForm(
            @JsonProperty("name") String title,
            @JsonProperty("series") String series,
            @JsonProperty("internal_id") int internalId
    ) {}

    @Data @NoArgsConstructor
    public static class TestForm {
        private String title;
        private String series;
        private int internalId;
    }

    private final ComicsRepository repository;

    @PostMapping("/track")
    public String track(TestForm form) {
        var comic = repository.findByInternalId(form.getInternalId());
        if (comic.isPresent()) {
            repository.delete(comic.get());
            return "follow";
        }

        repository.save(new Comic(form.getInternalId(), form.getTitle(), form.getSeries()));
        return "unfollow";
    }
}
