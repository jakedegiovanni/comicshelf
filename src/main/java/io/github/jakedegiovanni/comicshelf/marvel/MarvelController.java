package io.github.jakedegiovanni.comicshelf.marvel;

import io.github.jakedegiovanni.comicshelf.views.Page;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Controller;
import org.springframework.ui.Model;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;

import java.io.IOException;
import java.time.Clock;
import java.time.LocalDate;

@Controller
@RequestMapping("/marvel-unlimited/comics")
@RequiredArgsConstructor
public class MarvelController {

    private final Clock clock;
    private final MarvelClient client;

    @GetMapping
    public String weeklyComics(Model model) throws IOException, InterruptedException {
        var now = LocalDate.now(clock);
        Page.setupModel(model, client.weeklyComics(now), now);
        return "marvel-unlimited/comics";
    }
}
