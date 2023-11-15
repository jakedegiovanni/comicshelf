package io.github.jakedegiovanni.comicshelf.marvel;

import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.ObjectMapper;
import io.github.jakedegiovanni.comicshelf.marvel.model.Comic;
import io.github.jakedegiovanni.comicshelf.marvel.model.DataWrapper;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;

import java.io.IOException;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.time.DayOfWeek;
import java.time.LocalDate;
import java.time.format.DateTimeFormatter;

@Service
@RequiredArgsConstructor
@Slf4j
public class MarvelClient {

    private record DateRange(LocalDate start, LocalDate end) {};

    private final HttpClient httpClient;
    private final ObjectMapper objectMapper;
    private final MarvelConfig config;

    // todo handle bad status codes

    public DataWrapper<Comic> weeklyComics(LocalDate today) throws IOException, InterruptedException {
        DateRange dates = getDateRange(today);
        var endpoint = String.format(
                "/comics?format=comic&formatType=comic&noVariants=true&dateRange=%s,%s&hasDigitalIssue=true&orderBy=issueNumber&limit=100",
                dates.start().format(DateTimeFormatter.ISO_LOCAL_DATE),
                dates.end().format(DateTimeFormatter.ISO_LOCAL_DATE)
        );
        HttpResponse<String> response = get(endpoint);

        log.debug("raw json: {}", response.body());
        return objectMapper.readValue(response.body(), new TypeReference<>() {
        });
    }

    private HttpResponse<String> get(String endpoint) throws IOException, InterruptedException {
        var uri = config.uri(endpoint);
        var request = HttpRequest.newBuilder()
                .GET()
                .uri(uri)
//                .setHeader("", "") todo etag cache
                .build();
        return httpClient.send(request, HttpResponse.BodyHandlers.ofString());
    }

    private DateRange getDateRange(LocalDate today) {
        var localDate = today.minusMonths(3);

        while (localDate.getDayOfWeek() != DayOfWeek.SUNDAY) {
            localDate = localDate.minusDays(1);
        }

        return new DateRange(localDate.minusDays(7), localDate.minusDays(1));
    }
}
