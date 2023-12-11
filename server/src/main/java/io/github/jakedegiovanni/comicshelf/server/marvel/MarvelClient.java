package io.github.jakedegiovanni.comicshelf.server.marvel;

import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.ObjectMapper;
import io.github.jakedegiovanni.comicshelf.server.marvel.model.Comic;
import io.github.jakedegiovanni.comicshelf.server.marvel.model.DataWrapper;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;

import java.io.IOException;
import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.security.MessageDigest;
import java.security.NoSuchAlgorithmException;
import java.time.Clock;
import java.time.DayOfWeek;
import java.time.Instant;
import java.time.LocalDate;
import java.util.HexFormat;
import java.util.Optional;
import java.util.concurrent.ConcurrentHashMap;

import static java.nio.charset.StandardCharsets.UTF_8;
import static java.time.format.DateTimeFormatter.ISO_LOCAL_DATE;

@Service
@RequiredArgsConstructor
@Slf4j
public class MarvelClient {

    private record DateRange(LocalDate start, LocalDate end) {}

    // todo - cache evictions
    private final ConcurrentHashMap<String, DataWrapper<Comic>> comicCache = new ConcurrentHashMap<>();

    private final HttpClient httpClient;
    private final ObjectMapper objectMapper;
    private final MarvelConfig config;
    private final Clock clock;

    public DataWrapper<Comic> weeklyComics(LocalDate today) throws IOException, InterruptedException {
        log.debug("getting weekly comics");

        DateRange dates = getDateRange(today);
        var endpoint = STR
                ."/comics?format=comic&formatType=comic&noVariants=true&dateRange=\{dates.start().format(ISO_LOCAL_DATE)},\{dates.end().format(ISO_LOCAL_DATE)}&hasDigitalIssue=true&orderBy=issueNumber&limit=100";
        return get(endpoint);
    }

    private DataWrapper<Comic> get(String endpoint) throws IOException, InterruptedException {
        var uri = URI.create(STR."\{config.getBaseUrl()}\{endpoint}&\{auth()}");
        var request = HttpRequest.newBuilder()
                .GET()
                .uri(uri);

        if (comicCache.containsKey(endpoint)) {
            var comicDataWrapper = comicCache.get(endpoint);
            var etag = comicDataWrapper.getEtag();
            log.debug("endpoint {} present in cache, using etag {}", endpoint, etag);
            request.setHeader("If-None-Match", comicDataWrapper.getEtag());
        }

        var resp = httpClient.send(request.build(), HttpResponse.BodyHandlers.ofString());
        // todo handle bad status codes
        if (resp.statusCode() == 304) {
            log.debug("response not modified, using etag cache entry");
            return Optional.ofNullable(comicCache.get(endpoint)).orElseThrow(() -> new RuntimeException("expected entry not found"));
        }

        var result = objectMapper.readValue(resp.body(), new TypeReference<DataWrapper<Comic>>() {});
        comicCache.put(endpoint, result);
        log.debug("response cached, returning");
        return result;
    }

    private String auth() {
        var now = Instant.now(clock).toEpochMilli();

        MessageDigest md;
        try {
            md = MessageDigest.getInstance("MD5");
        } catch (NoSuchAlgorithmException e) {
            log.error("could not get md5 message digest", e);
            // todo handle
            throw new RuntimeException(e);
        }

        String hash = HexFormat.of().formatHex(md.digest(
                STR."\{now}\{config.getPriv()}\{config.getPub()}".getBytes(UTF_8)
        ));
        return STR."ts=\{now}&hash=\{hash}&apikey=\{config.getPub()}";
    }

    private DateRange getDateRange(LocalDate today) {
        var localDate = today.minusMonths(3);

        while (localDate.getDayOfWeek() != DayOfWeek.SUNDAY) {
            localDate = localDate.minusDays(1);
        }

        return new DateRange(localDate.minusDays(7), localDate.minusDays(1));
    }
}
