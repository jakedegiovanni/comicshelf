package io.github.jakedegiovanni.comicshelf.server;

import com.fasterxml.jackson.databind.ObjectMapper;
import io.github.jakedegiovanni.comicshelf.sdk.marvel.MarvelClient;
import io.github.jakedegiovanni.comicshelf.sdk.marvel.MarvelConfig;
import lombok.Data;
import lombok.NoArgsConstructor;
import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.boot.context.properties.EnableConfigurationProperties;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.data.jpa.repository.config.EnableJpaRepositories;
import org.springframework.web.filter.CommonsRequestLoggingFilter;

import java.net.http.HttpClient;
import java.time.Clock;
import java.util.concurrent.Executors;

@Configuration
@EnableConfigurationProperties
@ConfigurationProperties(prefix = "comicshelf")
@Data @NoArgsConstructor
@EnableJpaRepositories
public class AppConfig {

    private MarvelConfig marvel;

    @Bean
    public MarvelClient marvelClient(
            HttpClient httpClient,
            ObjectMapper objectMapper,
            AppConfig config,
            Clock clock
    ) {
        return new MarvelClient(httpClient, objectMapper, config.getMarvel(), clock);
    }

    @Bean
    public HttpClient httpClient() {
        return HttpClient.newBuilder()
                .version(HttpClient.Version.HTTP_1_1)
                .executor(Executors.newVirtualThreadPerTaskExecutor())
                .build();
    }

    @Bean
    public Clock clock() {
        return Clock.systemUTC();
    }

    @Bean
    public CommonsRequestLoggingFilter commonsRequestLoggingFilter() {
        var filter = new CommonsRequestLoggingFilter();
        filter.setIncludeQueryString(true);
        filter.setIncludeHeaders(false);
        return filter;
    }
}
