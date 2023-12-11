package io.github.jakedegiovanni.comicshelf.server.marvel;

import io.github.jakedegiovanni.comicshelf.server.comics.ComicsRepository;
import io.github.jakedegiovanni.comicshelf.server.views.Page;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Controller;
import org.springframework.ui.Model;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;

import java.io.IOException;
import java.time.Clock;
import java.time.LocalDate;
import java.util.concurrent.Executors;

@Controller
@RequestMapping("/marvel-unlimited/comics")
@RequiredArgsConstructor
@Slf4j
public class MarvelController {

    private final Clock clock;
    private final MarvelClient client;
    private final ComicsRepository repository;

    @GetMapping
    public String weeklyComics(Model model) throws IOException, InterruptedException {
        var now = LocalDate.now(clock);
        var result = client.weeklyComics(now);

        try (var executor = Executors.newVirtualThreadPerTaskExecutor()) {
            result.getData().getResults().forEach(comic -> {
                executor.submit(() -> {
                    log.debug("checking if following: {}", comic.getTitle());
                    repository.findByInternalId(comic.getId()).ifPresent(ignore -> comic.setFollowing(true));
                });
            });
        }

        Page.setupModel(model, result, now);
        return "marvel-unlimited/comics";
    }
}
