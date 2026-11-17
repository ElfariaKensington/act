package artifacts

import (
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestStandaloneArtifactRouteAuthenticationPolicy(t *testing.T) {
	cases := []struct {
		method, path string
		protected    bool
	}{
		{http.MethodGet, "/", true},
		{http.MethodGet, "/artifacts", true},
		{http.MethodHead, "/artifacts", false},
		{http.MethodPost, "/artifacts", false},
		{http.MethodGet, "/local-artifacts/summary?runId=123", false},
		{http.MethodGet, "/local-artifacts/download?path=123/file.zip", false},
		{http.MethodGet, "/local-artifacts/icon/download.svg", false},
		{http.MethodGet, "/local-artifacts/name.svg?name=example", false},
		{http.MethodGet, "/download/123", false},
		{http.MethodGet, "/artifact/123/file.zip", false},
		{http.MethodGet, "/_apis/pipelines/workflows/123/artifacts", false},
		{http.MethodPost, "/_apis/pipelines/workflows/123/artifacts", false},
		{http.MethodPatch, "/_apis/pipelines/workflows/123/artifacts", false},
		{http.MethodPut, "/upload/123?itemPath=file.zip", false},
		{http.MethodPost, ArtifactV4RouteBase + "/ListArtifacts", false},
		{http.MethodPost, ArtifactV4RouteBase + "/GetSignedArtifactURL", false},
		{http.MethodPost, ArtifactV4RouteBase + "/DeleteArtifact", false},
		{http.MethodPost, ArtifactV4RouteBase + "/CreateArtifact", false},
		{http.MethodPost, ArtifactV4RouteBase + "/FinalizeArtifact", false},
		{http.MethodPut, ArtifactV4RouteBase + "/UploadArtifact", false},
		{http.MethodGet, ArtifactV4RouteBase + "/DownloadArtifact", false},
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
