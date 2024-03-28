<?php

namespace Comicshelf\Services;

use Comicshelf\Util\Defer;
use DateInterval;
use DateTimeImmutable;
use Exception;
use Phalcon\Logger\Logger;

class MarvelClient
{

    public static string $DATE_FORMAT = 'Y-m-d';

    private string $priv;
    private string $pub;
    private Logger $logger;

    public function __construct(Logger $logger)
    {
        $this->priv = file_get_contents('comicclient/marvel/priv.txt');
        $this->pub = file_get_contents('comicclient/marvel/pub.txt');
        $this->logger = $logger;
    }

    /**
     * @throws Exception
     */
    public function getWeeklyComics(DateTimeImmutable $date): mixed
    {
        $this->logger->info("Getting weekly comics for " . $date->format(self::$DATE_FORMAT));

        if (date('l', $date->getTimestamp()) == 'Sunday') {
            $date = $date->sub(DateInterval::createFromDateString('1 day'));
        }

        [$d1, $d2] = $this->weekRange($this->marvelUnlimitedDate($date));

        $d1 = $d1->format(self::$DATE_FORMAT);
        $d2 = $d2->format(self::$DATE_FORMAT);

        $endpoint = "/comics?format=comic&formatType=comic&noVariants=true&dateRange=$d1,$d2&hasDigitalIssue=true&orderBy=issueNumber&limit=100";
        return $this->performRequest($endpoint);
    }

    private function auth(): string
    {
        $ts = time();
        $ts = "$ts";
        $hash = md5("$ts$this->priv$this->pub");

        return "ts=$ts&hash=${hash}&apikey=$this->pub";
    }

    /**
     * @throws Exception
     */
    private function performRequest(string $endpoint): mixed
    {
        $auth = $this->auth();
        $curl = curl_init("https://gateway.marvel.com/v1/public$endpoint&$auth");
        $defer = new Defer(function () use ($curl) {
            curl_close($curl);
        });

        curl_setopt_array($curl, [
            CURLOPT_RETURNTRANSFER => true,
            CURLOPT_SSL_VERIFYPEER => false, // todo
        ]);

        $result = curl_exec($curl);
        if (!$result) {
            throw new Exception("failed curl exec: " . curl_error($curl));
        }

        $code = curl_getinfo($curl, CURLINFO_RESPONSE_CODE);
        if ($code != 200) {
            throw new Exception("non 200 status code: ${code}");
        }

        return json_decode($result);
    }

    /**
     * @param DateTimeImmutable $date
     * @return DateTimeImmutable
     */
    private function marvelUnlimitedDate(DateTimeImmutable $date): DateTimeImmutable
    {
        return $date->sub(DateInterval::createFromDateString('3 months'));
    }

    /**
     * @param DateTimeImmutable $date
     * @return DateTimeImmutable[]
     */
    private function weekRange(DateTimeImmutable $date): array
    {
        while (date('l', $date->getTimestamp()) != 'Sunday') {
            $date = $date->sub(DateInterval::createFromDateString('1 day'));
        }

        $date = $date->sub(DateInterval::createFromDateString('7 days'));

        return [
            $date,
            $date->add(DateInterval::createFromDateString('6 days')),
        ];
    }
}
