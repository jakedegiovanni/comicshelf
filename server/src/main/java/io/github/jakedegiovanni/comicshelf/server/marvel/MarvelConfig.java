package io.github.jakedegiovanni.comicshelf.server.marvel;

import lombok.Data;
import lombok.NoArgsConstructor;
import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.context.annotation.Configuration;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;

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
