package io.github.jakedegiovanni.comicshelf.marvel;

import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Controller;
import org.springframework.ui.Model;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;

import java.io.IOException;
import java.time.Clock;
import java.time.LocalDate;
import java.time.format.DateTimeFormatter;

@Controller
@RequestMapping("/marvel-unlimited/comics")
@RequiredArgsConstructor
public class MarvelController {

    private final Clock clock;
    private final MarvelClient client;

    @GetMapping
    public String weeklyComics(Model model) throws IOException, InterruptedException {
        var now = LocalDate.now(clock);
        model.addAttribute("model", client.weeklyComics(now));
        model.addAttribute("pageEndpoint", "/marvel-unlimited/comics");
        model.addAttribute("date", now.format(DateTimeFormatter.ISO_LOCAL_DATE));
        return "marvel-unlimited/comics";
    }
}
