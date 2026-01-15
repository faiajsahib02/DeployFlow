import { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { deploymentApi } from '@/services/api';
import type { Project, Deployment } from '@/services/api';
import {
  Loader2,
  Rocket,
  CheckCircle2,
  XCircle,
  ExternalLink,
  ArrowLeft,
} from 'lucide-react';

interface DeploymentWorkspaceProps {
  project: Project;
  onBack: () => void;
}

type DeployState = 'idle' | 'deploying' | 'success' | 'error';

const DEFAULT_CODE = `from flask import Flask

app = Flask(__name__)

@app.route('/')
def hello():
    return 'Hello from DeployFlow! 🚀'

@app.route('/predict')
def predict():
    return {'prediction': 42, 'confidence': 0.95}

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8000)
`;

export function DeploymentWorkspace({ project, onBack }: DeploymentWorkspaceProps) {
  const [code, setCode] = useState(DEFAULT_CODE);
  const [deployState, setDeployState] = useState<DeployState>('idle');
  const [deployment, setDeployment] = useState<Deployment | null>(null);
  const [error, setError] = useState<string | null>(null);

  const handleDeploy = async () => {
    setDeployState('deploying');
    setError(null);

    try {
      const result = await deploymentApi.create(project.id, code);
      setDeployment(result);
      setDeployState('success');
    } catch (err) {
      setDeployState('error');
      if (err instanceof Error) {
        setError(err.message);
      } else {
        setError('Deployment failed. Check the console for details.');
      }
    }
  };

  const deploymentUrl = `http://${project.name}.localhost:8000`;

  return (
    <div className="min-h-screen flex flex-col">
      {/* Header */}
      <header className="border-b px-6 py-4">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-4">
            <Button variant="ghost" size="sm" onClick={onBack} className="gap-2">
              <ArrowLeft className="h-4 w-4" />
              Back
            </Button>
            <div className="h-6 w-px bg-border" />
            <div>
              <h1 className="font-semibold">{project.name}</h1>
              <p className="text-xs text-muted-foreground font-mono">{project.id}</p>
            </div>
          </div>
          <div className="flex items-center gap-2">
            <Rocket className="h-5 w-5 text-primary" />
            <span className="font-semibold">DeployFlow</span>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <div className="flex-1 flex">
        {/* Left Panel - Code Editor */}
        <div className="flex-1 flex flex-col border-r">
          <div className="px-4 py-3 border-b bg-secondary/30">
            <span className="text-sm font-medium">app.py</span>
          </div>
          <div className="flex-1 p-4">
            <textarea
              value={code}
              onChange={(e) => setCode(e.target.value)}
              className="w-full h-full bg-[#0d1117] text-[#c9d1d9] font-mono text-sm p-4 rounded-lg border border-border resize-none focus:outline-none focus:ring-2 focus:ring-ring"
              spellCheck={false}
              placeholder="Paste your Python/Flask code here..."
            />
          </div>
          <div className="p-4 border-t">
            <Button
              size="lg"
              onClick={handleDeploy}
              disabled={deployState === 'deploying'}
              className="w-full gap-2"
            >
              {deployState === 'deploying' ? (
                <>
                  <Loader2 className="h-5 w-5 animate-spin" />
                  Deploying...
                </>
              ) : (
                <>
                  <Rocket className="h-5 w-5" />
                  Deploy to Production
                </>
              )}
            </Button>
          </div>
        </div>

        {/* Right Panel - Status */}
        <div className="w-[400px] flex flex-col p-6 space-y-6">
          <Card>
            <CardHeader>
              <CardTitle className="text-base">Deployment Status</CardTitle>
            </CardHeader>
            <CardContent>
              {deployState === 'idle' && (
                <div className="text-center py-8 text-muted-foreground">
                  <p>Ready to deploy.</p>
                  <p className="text-sm mt-2">
                    Paste your Flask code and click Deploy.
                  </p>
                </div>
              )}

              {deployState === 'deploying' && (
                <div className="text-center py-8 space-y-4">
                  <div className="flex justify-center">
                    <div className="relative">
                      <div className="h-16 w-16 rounded-full border-4 border-primary/20" />
                      <div className="absolute inset-0 h-16 w-16 rounded-full border-4 border-primary border-t-transparent animate-spin" />
                    </div>
                  </div>
                  <div className="space-y-1">
                    <p className="font-medium">Building & Deploying...</p>
                    <p className="text-sm text-muted-foreground">
                      This may take 20-30 seconds
                    </p>
                  </div>
                  <div className="text-xs text-muted-foreground space-y-1">
                    <p>📦 Building Docker Image...</p>
                    <p>🚀 Starting Container...</p>
                  </div>
                </div>
              )}

              {deployState === 'success' && deployment && (
                <div className="text-center py-6 space-y-4">
                  <CheckCircle2 className="h-16 w-16 text-green-500 mx-auto" />
                  <div>
                    <p className="font-medium text-green-500">
                      Deployment Successful!
                    </p>
                    <p className="text-sm text-muted-foreground mt-1">
                      Your app is now live
                    </p>
                  </div>
                  <div className="pt-4 space-y-3">
                    <div>
                      <p className="text-xs text-muted-foreground mb-1">
                        Your App URL
                      </p>
                      <a
                        href={deploymentUrl}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="inline-flex items-center gap-2 text-primary hover:underline font-mono text-sm"
                      >
                        {deploymentUrl}
                        <ExternalLink className="h-3 w-3" />
                      </a>
                    </div>
                    <div className="text-xs text-muted-foreground">
                      <p>
                        Port: <span className="font-mono">{deployment.port}</span>
                      </p>
                      <p>
                        Container:{' '}
                        <span className="font-mono">
                          {deployment.container_id?.slice(0, 12)}
                        </span>
                      </p>
                    </div>
                  </div>
                </div>
              )}

              {deployState === 'error' && (
                <div className="text-center py-6 space-y-4">
                  <XCircle className="h-16 w-16 text-red-500 mx-auto" />
                  <div>
                    <p className="font-medium text-red-500">Deployment Failed</p>
                    <p className="text-sm text-muted-foreground mt-2">{error}</p>
                  </div>
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => setDeployState('idle')}
                  >
                    Try Again
                  </Button>
                </div>
              )}
            </CardContent>
          </Card>

          {/* Project Info Card */}
          <Card>
            <CardHeader>
              <CardTitle className="text-base">Project Info</CardTitle>
            </CardHeader>
            <CardContent className="space-y-3 text-sm">
              <div>
                <p className="text-muted-foreground">Name</p>
                <p className="font-medium">{project.name}</p>
              </div>
              <div>
                <p className="text-muted-foreground">Project ID</p>
                <p className="font-mono text-xs">{project.id}</p>
              </div>
              <div>
                <p className="text-muted-foreground">Created</p>
                <p>{new Date(project.created_at).toLocaleString()}</p>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}
