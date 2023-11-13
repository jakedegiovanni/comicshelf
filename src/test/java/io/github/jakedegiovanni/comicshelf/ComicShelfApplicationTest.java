package io.github.jakedegiovanni.comicshelf;

import org.junit.jupiter.api.Test;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.test.context.ActiveProfiles;

@SpringBootTest(webEnvironment = SpringBootTest.WebEnvironment.RANDOM_PORT)
@ActiveProfiles({"local"})
class ComicShelfApplicationTest {

    @Test
    void applicationStarts() {}
}
