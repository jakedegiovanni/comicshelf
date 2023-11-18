package io.github.jakedegiovanni.comicshelf.marvel;

import io.github.jakedegiovanni.comicshelf.comics.ComicsRepository;
import io.github.jakedegiovanni.comicshelf.views.Page;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Controller;
import org.springframework.ui.Model;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;

import java.io.IOException;
import java.time.Clock;
import java.time.LocalDate;
import java.util.concurrent.CompletableFuture;
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

        // tried to use virtual threads here over the loop and, alternatively, as the executor pool in completable future
        // but something within the database layers was causing the thread to pin.
        // will likely be able to migrate to virtual threads for the "following" population
        // in either future library versions or future jdk releases
        CompletableFuture.allOf(
                result.getData().getResults().stream()
                        .map(CompletableFuture::completedFuture)
                        .map(f -> f.thenAcceptAsync(c -> {
                            log.debug("checking if following: {}", c.getTitle());
                            repository.findByInternalId(c.getId()).ifPresent((c1) -> c.setFollowing(true));
                        }))
                        .toArray(CompletableFuture[]::new)
                )
                .join();

        Page.setupModel(model, result, now);
        return "marvel-unlimited/comics";
    }
}
