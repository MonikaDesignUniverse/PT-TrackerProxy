<?php

namespace App\Helpers;

use App\Models\User;
use Illuminate\Http\Request;

class PTTPHelper
{
	private static function ProcessPTTPListen(?string $pttpAddr, int|string|null $pttpPort): array {
		if ($pttpAddr == null || !filter_var($pttpAddr, FILTER_VALIDATE_IP, FILTER_FLAG_IPV4 | FILTER_FLAG_NO_PRIV_RANGE | FILTER_FLAG_NO_RES_RANGE)) {
			$pttpAddr = '127.0.0.1';
		}
		$pttpPort = null;
		if ($pttpPort == null || !is_numeric($pttpPort) || $pttpPort <= 1024 || $pttpPort >= 65535) {
			$pttpPort = 7887;
		}

		return array($pttpAddr, $pttpPort);
	}
	public static function CheckPTTPVersion(Request $request): ?string {
		$pttpVersion = $request->headers->get('x-pttp-version');
		if ($pttpVersion == null) {
			return null;
		}

		$pttpVersion = explode(' ', $pttpVersion, 2)[0];
		if (strlen($pttpVersion) > 24 || str_contains($pttpVersion, '(') || str_contains($pttpVersion, ')') || str_contains($pttpVersion, ',')) {
			return null;
		}

		return $pttpVersion;
	}
	public static function GetPTTPInfo(Request $request): array {
		$pttpVersion = self::CheckPTTPVersion($request);
		if ($pttpVersion !== null) {
			list($pttpAddr, $pttpPort) = self::ProcessPTTPListen($request->headers->get('x-pttp-listenaddr'), $request->headers->get('x-pttp-listenport'));
			return array('pttpVersion' => $pttpVersion, 'pttpAddr' => $pttpAddr, 'pttpPort' => $pttpPort);
		}

		return array('pttpVersion' => null);
	}
	public static function GetPTTPInfoFromUser(User $user): array {
		if (!empty($user->pttp_listen)) {
			$pttpListenSplit = explode(':', trim($user->pttp_listen), 3);
			if (count($pttpListenSplit) === 2) {
				list($pttpAddr, $pttpPort) = self::ProcessPTTPListen($pttpListenSplit[0], $pttpListenSplit[1]);
				return array('pttpVersion' => 'Unknown (FromUser)', 'pttpAddr' => $pttpAddr, 'pttpPort' => $pttpPort);
			}
		}

		return array('pttpVersion' => null);
	}
	public static function GetPTTPURL(array $pttpInfo, string $url) {
		if ($pttpInfo['pttpVersion'] !== null) {
			return str_replace(array('http://', 'https://'), "http://{$pttpInfo['pttpAddr']}:{$pttpInfo['pttpPort']}/", $url);
		}

		return $url;
	}
	public static function GetPTTPIPs(Request $request): array {
		$pttpIPs = array();
		$pttpVersion = self::CheckPTTPVersion($request);
		if ($pttpVersion !== null) {
			$pttpIP4 = $request->headers->get('x-pttp-ip4');
			$pttpIP6 = $request->headers->get('x-pttp-ip6');
			if ($pttpIP4 != null && filter_var($pttpIP4, FILTER_VALIDATE_IP, FILTER_FLAG_IPV4 | FILTER_FLAG_NO_PRIV_RANGE | FILTER_FLAG_NO_RES_RANGE)) {
				array_push($pttpIPs, $pttpIP4);
			}
			if ($pttpIP6 != null && filter_var($pttpIP6, FILTER_VALIDATE_IP, FILTER_FLAG_IPV6 | FILTER_FLAG_NO_PRIV_RANGE | FILTER_FLAG_NO_RES_RANGE)) {
				array_push($pttpIPs, $pttpIP6);
			}
		}

		return array('pttpVersion' => $pttpVersion, 'pttpIPs' => $pttpIPs);
	}
}

/*
Part RSS.
$pttpInfo = PTTPHelper::GetPTTPInfo($request);
$rssURL = PTTPHelper::GetPTTPURL($pttpInfo, route('rss.show.rsskey', ['id' => $rss->id, 'rsskey' => $user->rsskey]));
$downloadURL = PTTPHelper::GetPTTPURL($pttpInfo, route('torrent.download.rsskey', ['id' => $data->id, 'rsskey' => $user->rsskey ]));
*/

/*
Part TorrentDownload.
if (!$usePTTPFromUser) {
    $pttpInfo = PTTPHelper::GetPTTPInfo($request);
} else { // If user downloads torrent directly from site (not via PTTP).
    $pttpInfo = PTTPHelper::GetPTTPInfoFromUser($user);
}
$dict['announce'] = PTTPHelper::GetPTTPURL($pttpInfo, \route('announce', ['passkey' => $user->passkey]));
*/

/*
Part Announce.
$ip = $request->getClientIp();
$userAgent = $request->headers->get('user-agent');
$processLoop = true;
$processPTTP = false;
$pttpIPs = PTTPHelper::GetPTTPIPs($request);
if ($pttpIPs['pttpVersion'] !== null) {
	$userAgent .= " PT-TrackerProxy/{$pttpIPs['pttpVersion']}";
}

while ($processLoop) {
	if ($pttpIPs['pttpVersion'] !== null && ($pttpIP = array_shift($pttpIPs['pttpIPs'])) !== null) {
		$processPTTP = true;
		$ip = $pttpIP;
	} else {
		$processLoop = false;
		if ($processPTTP) {
			break;
		}
	}
	ProcessAnnounce(...);
}
*/
