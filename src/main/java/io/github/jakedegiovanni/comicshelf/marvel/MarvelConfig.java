package io.github.jakedegiovanni.comicshelf.marvel;

import lombok.Data;
import lombok.NoArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.context.annotation.Configuration;

import java.io.IOException;
import java.net.URI;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.security.MessageDigest;
import java.security.NoSuchAlgorithmException;
import java.time.Clock;
import java.time.Instant;

@Configuration
@Slf4j
@ConfigurationProperties(prefix = "marvel")
@Data @NoArgsConstructor
class MarvelConfig {

    private String url;

    URI uri(String endpoint) {
        String pub;
        String priv;
        try {
            pub = Files.readString(Path.of("pub.txt"));
            priv = Files.readString(Path.of("priv.txt"));
        } catch (IOException e) {
            log.error("could not read auth files", e);
            // todo handle
            throw new RuntimeException(e);
        }

        var now = Instant.now(Clock.systemUTC()).toEpochMilli();
        String hash;
        try {
            var digest = MessageDigest.getInstance("MD5")
                    .digest(
                            String.format(
                                    "%d%s%s",
                                    now,
                                    priv,
                                    pub
                            ).getBytes(StandardCharsets.UTF_8)
                    );
            var sb = new StringBuilder();
            for (var b : digest) {
                sb.append(String.format("%02x", b));
            }

            hash = sb.toString();
        } catch (NoSuchAlgorithmException e) {
            log.error("could not get md5 message digest", e);
            // todo handle
            throw new RuntimeException(e);
        }

        return URI.create(String.format(
                "%s%s&ts=%d&hash=%s&apikey=%s",
                url,
                endpoint,
                now,
                hash,
                pub
        ));
    }
}
