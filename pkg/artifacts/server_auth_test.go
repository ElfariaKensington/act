package artifacts

import (
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestStandaloneArtifactRouteAuthenticationPolicy(t *testing.T) {
    cases := []struct {
        method, path string
        protected bool
    }{
        {"GET", "/", true},
        {"GET", "/artifacts", true},
        {"GET", "/local-artifacts/summary?runId=123", true},
        {"GET", "/local-artifacts/download?path=123/file.zip", true},
        {"GET", "/local-artifacts/icon/download.svg", true},
        {"GET", "/download/123", true},
        {"GET", "/artifact/123/file.zip", true},
        {"GET", "/_apis/pipelines/workflows/123/artifacts", true},
        {"POST", ArtifactV4RouteBase + "/ListArtifacts", true},
        {"POST", ArtifactV4RouteBase + "/GetSignedArtifactURL", true},
        {"POST", ArtifactV4RouteBase + "/DeleteArtifact", true},
        {"POST", ArtifactV4RouteBase + "/CreateArtifact", false},
        {"POST", ArtifactV4RouteBase + "/FinalizeArtifact", false},
        {"PUT", ArtifactV4RouteBase + "/UploadArtifact", false},
        {"GET", ArtifactV4RouteBase + "/DownloadArtifact", false},
    }
    for _, tc := range cases {
        r := httptest.NewRequest(tc.method, tc.path, nil)
        if got := artifactRouteRequiresPassword(r); got != tc.protected {
            t.Errorf("%s %s: password required = %v, want %v", tc.method, tc.path, got, tc.protected)
        }
    }
}

func TestStandaloneArtifactBasicAuthentication(t *testing.T) {
    const password = "long-secret-password"
    for _, tc := range []struct {
        name, user, pass string
        authorized bool
    }{
        {"missing", "", "", false},
        {"incorrect", "admin", "wrong", false},
        {"correct", "admin", password, true},
    } {
        t.Run(tc.name, func(t *testing.T) {
            r := httptest.NewRequest(http.MethodGet, "/artifacts", nil)
            if tc.user != "" {
                r.SetBasicAuth(tc.user, tc.pass)
            }
            if got := artifactPasswordValid(r, password); got != tc.authorized {
                t.Fatalf("authorized = %v, want %v", got, tc.authorized)
            }
        })
    }
}
