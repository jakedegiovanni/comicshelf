package io.github.jakedegiovanni.comicshelf;

import org.springframework.boot.context.properties.EnableConfigurationProperties;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

import java.net.http.HttpClient;
import java.time.Clock;

@Configuration
@EnableConfigurationProperties
public class AppConfig {

    @Bean
    HttpClient httpClient() {
        return HttpClient.newBuilder()
                .version(HttpClient.Version.HTTP_1_1)
                .build();
    }

    @Bean
    Clock clock() {
        return Clock.systemUTC();
    }
}
