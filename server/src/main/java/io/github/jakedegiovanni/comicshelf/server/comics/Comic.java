package io.github.jakedegiovanni.comicshelf.server.comics;

import jakarta.persistence.*;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.util.UUID;

@Entity
@Table(name = "comic", indexes = {
        @Index(columnList = "internal_id")
})
@Data @NoArgsConstructor
public class Comic {

    @Id @GeneratedValue(strategy = GenerationType.UUID)
    private UUID id;

    @Column(name = "internal_id") private int internalId;
    private String title;
    private String series;

    public Comic(int internalId, String title, String series) {
        this.internalId = internalId;
        this.title = title;
        this.series = series;
    }
}
