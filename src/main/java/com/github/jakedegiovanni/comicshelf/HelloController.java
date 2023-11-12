package com.github.jakedegiovanni.comicshelf;

import org.springframework.stereotype.Controller;
import org.springframework.ui.Model;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;

@Controller
@RequestMapping("/hello")
public class HelloController {

    public record HelloModel(String msg) {};

    @GetMapping
    public String hello(Model model) {
        model.addAttribute("model", new HelloModel("world"));
        return "hello";
    }
}
