Import("env")
from SCons.Script import COMMAND_LINE_TARGETS

def buildWeb(source, target, env):
    env.Execute("cd web; pnpm build")
    print("Successfully built webui")

def convertImages(source, target, env):
    env.Execute("cd converter; go run .")
    print("Successfully converted images")

env.AddPreAction("buildprog", convertImages)
env.AddPreAction("uploadfs", buildWeb)