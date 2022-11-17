#http://www.pixiv.net/member_illust.php?mode=medium&illust_id=
#http://www.pixiv.net/member.php?id=
#https://www.pixiv.net/users/   
#https://www.pixiv.net/artworks/
param (
    [string[]]$id,
    [string[]]$uid
)
BEGIN {
    $urls = @()

    foreach ($i in $id ) {
        $urls += "https://www.pixiv.net/artworks/" + $i
    }
    foreach ($u in $uid ) {
        $urls += "https://www.pixiv.net/users/" + $u
    }

}
PROCESS {
    foreach ($url in $urls) {
        Start-Process msedge $url
    }
}
END { 
}

