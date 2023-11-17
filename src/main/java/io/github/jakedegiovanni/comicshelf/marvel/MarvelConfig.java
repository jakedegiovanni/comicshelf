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
@ConfigurationProperties(prefix = "marvel")
@Data @NoArgsConstructor
class MarvelConfig {

    private String baseUrl;
    private String pub;
    private String priv;

    void setPub(Path pub) throws IOException {
        this.pub = Files.readString(pub);
    }

    void setPriv(Path priv) throws IOException {
        this.priv = Files.readString(priv);
    }
}
