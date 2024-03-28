<?php

namespace Comicshelf\Controllers;

use Comicshelf\Services\MarvelClient;
use DateTimeImmutable;
use Exception;
use Phalcon\Logger\Logger;
use Phalcon\Mvc\Controller;

class ComicsController extends Controller
{

    private MarvelClient $client;
    private Logger $logger;

    public function initialize(): void
    {
        $this->client = $this->container['marvelClient'];
        $this->logger = $this->container['logger'];
    }

    /**
     * @throws Exception
     */
    public function indexAction(): void
    {
        if ($this->request->getQuery("date", "string", "") == "") {
            $this->logger->info("request doesn't have date parameter, redirecting");

            $date = date_create()->format(MarvelClient::$DATE_FORMAT);
            $this->response->redirect("/comics?date=$date");
            return;
        }

        $this->logger->info("Getting comics released");

        $this->view->title = "Comics";

        $date = $this->request->getQuery("date", "string");
        $this->view->date = $date;

        $date = DateTimeImmutable::createFromFormat(MarvelClient::$DATE_FORMAT, $date);

        $resp = $this->client->getWeeklyComics($date);
        $this->view->resp = $resp;
    }
}
