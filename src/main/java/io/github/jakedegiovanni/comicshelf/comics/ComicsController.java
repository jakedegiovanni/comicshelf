package io.github.jakedegiovanni.comicshelf.comics;

import jakarta.transaction.Transactional;
import lombok.Data;
import lombok.NoArgsConstructor;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestMapping;

@Controller
@RequestMapping("/api/comics")
@RequiredArgsConstructor
public class ComicsController {

    @Data @NoArgsConstructor
    public static class ComicForm {
        private String title;
        private String series;
        private int internalId;
    }

    private final ComicsRepository repository;

    @PostMapping("/follow")
    @Transactional
    public String follow(ComicForm form) {
        repository.findByInternalId(form.getInternalId()).ifPresentOrElse(
                comic -> {
                    comic.setSeries(form.getSeries());
                    comic.setTitle(form.getTitle());
                    repository.save(comic);
                },
                () -> repository.save(new Comic(form.getInternalId(), form.getTitle(), form.getSeries()))
        );
        return "unfollow";
    }

    @PostMapping("/unfollow")
    @Transactional
    public String unfollow(ComicForm form) {
        repository.deleteByInternalId(form.getInternalId());
        return "follow";
    }
}
