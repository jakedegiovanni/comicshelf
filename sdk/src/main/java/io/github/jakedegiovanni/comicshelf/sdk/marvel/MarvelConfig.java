package io.github.jakedegiovanni.comicshelf.sdk.marvel;

import lombok.Data;
import lombok.NoArgsConstructor;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;

@Data @NoArgsConstructor
public class MarvelConfig {

    private String baseUrl;
    private String pub;
    private String priv;

    public void setPub(Path pub) throws IOException {
        this.pub = Files.readString(pub);
    }

    public void setPriv(Path priv) throws IOException {
        this.priv = Files.readString(priv);
    }
}
